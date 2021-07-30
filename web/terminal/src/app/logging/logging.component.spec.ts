import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LoggingComponent } from './logging.component';

describe('LoggingComponent', () => {
  let component: LoggingComponent;
  let fixture: ComponentFixture<LoggingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ LoggingComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(LoggingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
