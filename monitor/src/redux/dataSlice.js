import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    data:[],
    deviceLoaded:false,
    device:null
}

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
      addDataItem:(state,action)=>{
        state.data=[...state.data,action.payload];
      },
      setDevice:(state,action)=>{
        console.log("setDevice:",action.payload);
        state.device=action.payload;
        state.deviceLoaded=true;
      }
    }
});

export const { 
  addDataItem,
  setDevice
} = dataSlice.actions

export default dataSlice.reducer