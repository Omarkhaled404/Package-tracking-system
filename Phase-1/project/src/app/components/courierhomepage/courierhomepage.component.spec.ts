import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CourierhomepageComponent } from './courierhomepage.component';

describe('CourierhomepageComponent', () => {
  let component: CourierhomepageComponent;
  let fixture: ComponentFixture<CourierhomepageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CourierhomepageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(CourierhomepageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
