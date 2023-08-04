import { configureStore } from '@reduxjs/toolkit'

import mqttReducer from './mqttSlice';
import dataReducer from './dataSlice';

export default configureStore({
  reducer: {
    mqtt:mqttReducer,
    data:dataReducer
  }
});