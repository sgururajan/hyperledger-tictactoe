
import { fakeAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatSidenavModule } from '@angular/material/sidenav';
import { TopnavbarComponent } from './topnavbar.component';

describe('TopnavbarComponent', () => {
  let component: TopnavbarComponent;
  let fixture: ComponentFixture<TopnavbarComponent>;

  beforeEach(fakeAsync(() => {
    TestBed.configureTestingModule({
      imports: [MatSidenavModule],
      declarations: [TopnavbarComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TopnavbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should compile', () => {
    expect(component).toBeTruthy();
  });
});
