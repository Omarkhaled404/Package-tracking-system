import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OwnerhomepageComponent } from './ownerhomepage.component';

describe('OwnerhomepageComponent', () => {
  let component: OwnerhomepageComponent;
  let fixture: ComponentFixture<OwnerhomepageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OwnerhomepageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(OwnerhomepageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
