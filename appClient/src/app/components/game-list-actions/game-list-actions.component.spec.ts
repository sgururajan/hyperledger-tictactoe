import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GameListActionsComponent } from './game-list-actions.component';

describe('GameListActionsComponent', () => {
  let component: GameListActionsComponent;
  let fixture: ComponentFixture<GameListActionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GameListActionsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GameListActionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
