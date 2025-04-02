package vault

type VaultOption func(*vault)

func WithWAL(wal WAL) VaultOption {
	return func(vault *vault) {
		vault.wal = wal
	}
}
