export function isLocalEnvironment(): boolean {
  const environment = (import.meta.env.VITE_APP_ENV || '').toUpperCase()

  return environment === 'LOCAL' || import.meta.env.DEV
}
