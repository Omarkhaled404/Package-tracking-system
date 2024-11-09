import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { adminorder, order } from '../../interfaces/order';
import {  OnInit } from '@angular/core';
import { OrdersService } from '../../ordersercices/orders.service';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms'; 
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
//import { BrowserModule } from '@angular/platform-browser';
@Component({
  selector: 'app-adminhomepage',
  standalone: true,
  imports: [CommonModule,FormsModule,CardModule,ButtonModule],
  templateUrl: './adminhomepage.component.html',
  styleUrl: './adminhomepage.component.css'
})
export class AdminhomepageComponent  {
  
  orders: adminorder[] = [];
  constructor(private ordersService: OrdersService, private router: Router) {}
  courierId: number | undefined; 
  
  getAllOrders(): void {
    this.ordersService.getAllOrders().subscribe({
      next: (orders: adminorder[]) => {
        this.orders = orders;  // Assign the fetched orders to the orders array
        console.log('Orders fetched:', this.orders);
      },
      error: (err: any) => {
        console.error('Error fetching orders:', err);
      }
    });
  }

  // goToOrderDetails(orderId: number): void {
  //   this.router.navigate(['adminmodifications', orderId]);
  // }
 
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

  deleteOrder(orderId: number): void {
    this.ordersService.deleteOrderbyadmin(orderId).subscribe({
      next: () => {
        this.orders = this.orders.filter(order => order.order_id !== orderId);
        console.log(`Order ${orderId} deleted successfully.`);
      },
      error: (err: any) => {
        console.error('Error deleting order:', err);
      }
    });
  }
  assignToCourier(orderId: number, courierId: number): void {
    if (courierId) {
      this.ordersService.assign(orderId, courierId).subscribe(
        (assignedOrder) => {
          console.log('Order assigned to courier:', assignedOrder);
          // Optionally update local order data after assignment
        },
        (error) => {
          console.error('Error assigning order to courier:', error);
        }
      );
    } else {
      console.log('Please enter a valid Courier ID.');
    }
  }
   
}