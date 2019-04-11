// Libraries
import {produce} from 'immer'

// Constants
import {
  AUTOREFRESH_DEFAULT_INTERVAL,
  AUTOREFRESH_DEFAULT_STATUS,
} from 'src/shared/constants'

// Types
import {Action, AutoRefreshStatus} from 'src/shared/actions/autorefresh'

export interface AutoRefreshState {
  status: AutoRefreshStatus
  interval: number
}

const initialState = (): AutoRefreshState => ({
  status: AUTOREFRESH_DEFAULT_STATUS,
  interval: AUTOREFRESH_DEFAULT_INTERVAL,
})

export const autoRefreshReducer = (state = initialState(), action: Action) =>
  produce(state, draftState => {
    switch (action.type) {
      case 'SET_AUTO_REFRESH_INTERVAL': {
        const {milliseconds} = action.payload

        if (milliseconds === 0){
          draftState.status = AutoRefreshStatus.Paused
        }

        draftState.interval = milliseconds

        return
      }

      case 'SET_AUTO_REFRESH_STATUS': {
        const {status} = action.payload

        draftState.status = status

        return
      }
    }
  })
