import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    data:[]
}

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
      addDataItem:(state,action)=>{
        state.data=[...state.data,action.payload];
      }
    }
});

export const { 
  addDataItem
} = dataSlice.actions

export default dataSlice.reducer