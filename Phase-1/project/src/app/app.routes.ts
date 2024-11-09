import { Routes } from '@angular/router';
import { LoginComponent } from './components/login/login.component';
import { RegisterComponent } from './components/register/register.component';

import { OwnerhomepageComponent } from './components/ownerhomepage/ownerhomepage.component';
// import { AdminhomepageComponent } from './components/adminhomepage/adminhomepage.component';
// import { CourierhomepageComponent } from './components/courierhomepage/courierhomepage.component';
 import { CreateOrderComponent } from './components/create-order/create-order.component';
import { AdminhomepageComponent } from './components/adminhomepage/adminhomepage.component';
import { CourierhomepageComponent } from './components/courierhomepage/courierhomepage.component';
// import { OrderDetailsComponent } from './components/order-details/order-details.component';

export const routes: Routes = [
    {
        path: 'login',
        component: LoginComponent
    },
    {
        path: 'register',
        component: RegisterComponent
    },
    {
        path: 'ownerhomepage',
        component: OwnerhomepageComponent
    },
    {
        path: 'courierhomepage',
        component: CourierhomepageComponent
    },
    {
        path: 'adminhomepage',
        component: AdminhomepageComponent
    },
     {
        path: 'create-order',
        component: CreateOrderComponent
    },
    // {
    //     path: 'order-details',
    //     component: OrderDetailsComponent
    // },
    {
        path: '',redirectTo: 'login', pathMatch: 'full'
    }
    
];