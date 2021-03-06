package keys

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/KiraCore/cosmos-sdk/client/flags"
	"github.com/KiraCore/cosmos-sdk/client/input"
	"github.com/KiraCore/cosmos-sdk/crypto/keyring"
	sdk "github.com/KiraCore/cosmos-sdk/types"
)

// migratePassphrase is used as a no-op migration key passphrase as a passphrase
// is not needed for importing into the Keyring keystore.
const migratePassphrase = "NOOP_PASSPHRASE"

// MigrateCommand migrates key information from legacy keybase to OS secret store.
func MigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate keys from the legacy (db-based) Keybase",
		Long: `Migrate key information from the legacy (db-based) Keybase to the new keyring-based Keybase.
For each key material entry, the command will prompt if the key should be skipped or not. If the key
is not to be skipped, the passphrase must be entered. The key will only be migrated if the passphrase
is correct. Otherwise, the command will exit and migration must be repeated.

It is recommended to run in 'dry-run' mode first to verify all key migration material.
`,
		Args: cobra.ExactArgs(0),
		RunE: runMigrateCmd,
	}

	cmd.Flags().Bool(flags.FlagDryRun, false, "Run migration without actually persisting any changes to the new Keybase")
	return cmd
}

func runMigrateCmd(cmd *cobra.Command, args []string) error {
	rootDir, _ := cmd.Flags().GetString(flags.FlagHome)

	// instantiate legacy keybase
	var legacyKb keyring.LegacyKeybase
	legacyKb, err := NewLegacyKeyBaseFromDir(rootDir)
	if err != nil {
		return err
	}

	defer legacyKb.Close()

	// fetch list of keys from legacy keybase
	oldKeys, err := legacyKb.List()
	if err != nil {
		return err
	}

	buf := bufio.NewReader(cmd.InOrStdin())
	keyringServiceName := sdk.KeyringServiceName()

	var (
		tmpDir   string
		migrator keyring.InfoImporter
	)

	if dryRun, _ := cmd.Flags().GetBool(flags.FlagDryRun); dryRun {
		tmpDir, err = ioutil.TempDir("", "migrator-migrate-dryrun")
		if err != nil {
			return errors.Wrap(err, "failed to create temporary directory for dryrun migration")
		}

		defer os.RemoveAll(tmpDir)

		migrator, err = keyring.NewInfoImporter(keyringServiceName, "test", tmpDir, buf)
	} else {
		backend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
		migrator, err = keyring.NewInfoImporter(keyringServiceName, backend, rootDir, buf)
	}

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(
			"failed to initialize keybase for service %s at directory %s",
			keyringServiceName, rootDir,
		))
	}

	for _, key := range oldKeys {
		legKeyInfo, err := legacyKb.Export(key.GetName())
		if err != nil {
			return err
		}

		keyName := key.GetName()
		keyType := key.GetType()

		cmd.PrintErrf("Migrating key: '%s (%s)' ...\n", key.GetName(), keyType)

		// allow user to skip migrating specific keys
		ok, err := input.GetConfirmation("Skip key migration?", buf, cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		if keyType != keyring.TypeLocal {
			if err := migrator.Import(keyName, legKeyInfo); err != nil {
				return err
			}

			continue
		}

		password, err := input.GetPassword("Enter passphrase to decrypt key:", buf)
		if err != nil {
			return err
		}

		// NOTE: A passphrase is not actually needed here as when the key information
		// is imported into the Keyring-based Keybase it only needs the password
		// (see: writeLocalKey).
		armoredPriv, err := legacyKb.ExportPrivKey(keyName, password, migratePassphrase)
		if err != nil {
			return err
		}

		if err := migrator.Import(keyName, armoredPriv); err != nil {
			return err
		}
	}

	return err
}
