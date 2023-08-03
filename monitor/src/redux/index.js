import { configureStore } from '@reduxjs/toolkit'

import mqttReducer from './frameSlice';

export default configureStore({
  reducer: {
    mqtt:mqttReducer
  }
});