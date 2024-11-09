import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { createOrder } from '../../interfaces/order';
import { OrdersService } from '../../ordersercices/orders.service';
import { AuthService } from '../../services/auth.service';
import { MessageService } from 'primeng/api';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';

@Component({
  selector: 'app-create-order',
  templateUrl: './create-order.component.html',
  styleUrls: ['./create-order.component.css'],
  standalone: true,
  imports: [FormsModule,ButtonModule,CardModule]
})
export class CreateOrderComponent {
  orderData: createOrder = {
    user_id: 6,
    pickup_location: '',
    dropoff_location: '',
    package_details: '',
    delivery_time: ''
  };

  constructor(
    private router: Router,
    private orderService: OrdersService,
    private authService: AuthService,
    private messageService: MessageService // Correctly inject MessageService
  ) {/*this.orderData.user_id = this.authService.getUserId() ?? 0;*/}

  submitOrder() {
    console.log('Order Data:', this.orderData);  // Add this line to log orderData
    this.orderService.createOrder(this.orderData).subscribe({
      next: (response) => {
        this.messageService.add({ severity: 'success', summary: 'Success', detail: 'Order created successfully' });
        console.log('Order creation response:', response);
      },
      error: (err) => {
        console.error('Order creation error:', err);
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Failed to create order' });
      }
    });
  }
  
}
