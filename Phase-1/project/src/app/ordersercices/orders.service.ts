import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {adminorder, createOrder, order} from '../interfaces/order';
import { statSync } from 'node:fs';

@Injectable({
  providedIn: 'root'
})
export class OrdersService {
  private baseURL="http://localhost:8000";
  constructor(private http: HttpClient) { }
  createOrder(orderData: createOrder): Observable<any> {
    return this.http.post<any>(`${this.baseURL}/orders`, orderData);
  }
  
  getUserOrders(userId: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseURL}/user/orders?user_id=${userId}`);
  }
  getOrderDetails(orderId: number): Observable<any> {
    return this.http.get<any>(`${this.baseURL}/order/details?order_id=${orderId}`);
  }
  getAllOrders(): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseURL}/admin/orders`);
  }
  deleteOrderbyadmin(orderId: number): Observable<void> {
    const url = `${this.baseURL}/admin/orders/delete?order_id=${orderId}`;
    return this.http.delete<void>(url);
  }
  getOrderById(orderId: number): Observable<any> {
    return this.http.get<any>(`${this.baseURL}/${orderId}`);
  }
  updateOrderStatus(orderId: number, status: string): Observable<order> {
    const url = `${this.baseURL}/admin/orders/update?order_id=${orderId}`;
    const body = { status };

    return this.http.put<order>(url, body);
  }
  assign(orderId: number, courierId: number){
    const url = `${this.baseURL}/order/${orderId}/assign/${courierId}`;  // URL with parameters
    return this.http.post<order>(url,null);
  }
  getCourierOrders(courierID: number | null): Observable<adminorder[]> {
    return this.http.get<adminorder[]>(`${this.baseURL}/courier/${courierID}/orders`);
  }
  declineorderByCourier(orderId: number, courierId: number | null): Observable<void> {
    const url = `${this.baseURL}/order/${orderId}/decline/${courierId}`;
    return this.http.delete<void>(url);
  }
}