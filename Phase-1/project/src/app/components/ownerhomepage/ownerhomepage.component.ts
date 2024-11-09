import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { OrdersService } from '../../ordersercices/orders.service'; // Updated import
import { order } from '../../interfaces/order';
import { CommonModule } from '@angular/common'; 
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';

@Component({
  selector: 'app-ownerhomepage',
  standalone: true,
  imports: [CommonModule,ButtonModule,CardModule],
  templateUrl: './ownerhomepage.component.html',
  styleUrls: ['./ownerhomepage.component.css']
})
export class OwnerhomepageComponent{
  userOrders: order[] = [];

  constructor(
    private router: Router,
    private authService: AuthService,
    private ordersService: OrdersService // Inject OrderService
  ) {}

  getOrders(): void {
    const userId = this.authService.getUserId();
    if (userId !== null) {
      this.ordersService.getUserOrders(userId).subscribe({
        next: (orders: any[]) => {
          this.userOrders = orders;
          console.log('User orders:', orders);
        },
        error: (err: any) => {
          console.error('Error fetching orders:', err);
        }
      });
    } else {
      console.error('No user ID found. User might not be logged in.');
    }
  }

  goToCreateOrder() {
    this.router.navigate(['create-order']);
  }

  cancelOrder(orderId: number, status: string): void {
    console.log('Order status:', status);  // Log the status to inspect its value
    if (status === 'pending') {
      this.ordersService.deleteOrderbyadmin(orderId).subscribe({
        next: () => {
          console.log('Order cancelled successfully');
          // Optionally, refresh the orders list or remove the canceled order from the UI
          this.userOrders = this.userOrders.filter(order => order.order_id !== orderId);
        },
        error: (err) => {
          console.error('Error cancelling order:', err);
        }
      });
    } else {
      console.log('Cannot cancel order. Status is not "Pending".');
    }
  }
  
  
  
}
