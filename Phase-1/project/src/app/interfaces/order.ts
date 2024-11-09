// src/app/models/order.interface.ts
export interface createOrder {
    user_id : number;
    pickup_location: string;
    dropoff_location: string;
    package_details: string;
    delivery_time: string;
}
export interface order {
    order_id: number;
    pickup_location: string;
    dropoff_location: string;
    package_details: string;
    delivery_time: string; // or Date if you’re using a Date object
    status: string; 
    courier_id?: number; // Optional field if courier info may or may not be available   
}
export interface adminorder{
    order_id: number;           
    user_id: number;           
    pickup_location: string;    
    dropoff_location: string;   
    package_details: string;    
    delivery_time: string;      
    status: string;             
    courier_id: number; 
  }