import { Component } from '@angular/core';
import { OrdersService } from '../../ordersercices/orders.service';
import { Router } from '@angular/router';
import { adminorder } from '../../interfaces/order';
//import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { AuthService } from '../../services/auth.service';
import { MessageService } from 'primeng/api'; 

@Component({
  selector: 'app-courierhomepage',
  standalone: true,
  imports: [FormsModule,CommonModule,CardModule,ButtonModule],
  templateUrl: './courierhomepage.component.html',
  styleUrl: './courierhomepage.component.css'
})
export class CourierhomepageComponent {
  message: string = '';

  acceptOrder(): void {
    //this.message = 'Order Accepted!';
    this.messageService.add({ severity: 'success', summary: 'Success', detail: 'Order Accepted' });
  }
  orders: adminorder[] = []; 
  constructor(private ordersService: OrdersService, private router: Router,private authService: AuthService,private messageService: MessageService) {}

  loadCourierOrders(): void {
    const courierId = this.authService.getUserId();
    this.ordersService.getCourierOrders(courierId).subscribe(
      (orders) => {
        this.orders = orders;
        console.log('Orders for courier:', orders);
      },
      (error) => {
        console.error('Error fetching courier orders:', error);
      }
    );
  }
  declineOrder(orderId: number): void {
    const courierId = this.authService.getUserId();
    this.ordersService.declineorderByCourier(orderId, courierId).subscribe(
      () => {
        console.log(`Order ${orderId} declined by courier ${courierId}`);
        // Remove the declined order from the local orders array or refresh the list
        this.orders = this.orders.filter(order => order.order_id !== orderId);
      },
      (error) => {
        console.error('Error declining the order:', error);
      }
    );
  }

  onStatusChange(order: adminorder): void {
    // Check if the status is not null or empty
    if (order.status) {
      this.updateOrderStatus(order.order_id, order.status);
    } else {
      console.log('Status is invalid. No update triggered.');
    }
  }
  updateOrderStatus(orderId: number, status: string) {
    if(status){
    this.ordersService.updateOrderStatus(orderId, status).subscribe(
      (updatedOrder) => {
        console.log(`Order ${orderId} status updated to: ${status}`, updatedOrder);
        // Optionally update the local orders array to reflect the new status
      },
      (error) => {
        console.error('Error updating order status:', error);
      }
    );
  }
  }
}
