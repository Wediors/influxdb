export enum AutoRefreshStatus {
  Active = 'active',
  Disabled = 'disabled',
  Paused = 'paused',
}

export type Action = SetAutoRefresh | SetAutoRefreshStatus

interface SetAutoRefresh {
  type: 'SET_AUTO_REFRESH_INTERVAL'
  payload: {milliseconds: number}
}

export const setAutoRefreshInterval = (
  milliseconds: number
): SetAutoRefresh => ({
  type: 'SET_AUTO_REFRESH_INTERVAL',
  payload: {milliseconds},
})

interface SetAutoRefreshStatus {
  type: 'SET_AUTO_REFRESH_STATUS'
  payload: {status: AutoRefreshStatus}
}

export const setAutoRefreshStatus = (
  status: AutoRefreshStatus
): SetAutoRefreshStatus => ({
  type: 'SET_AUTO_REFRESH_STATUS',
  payload: {status},
})
