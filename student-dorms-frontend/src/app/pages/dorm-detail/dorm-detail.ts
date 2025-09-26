import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Dorm , DormService  } from '../../services/dorm';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';


@Component({
  selector: 'app-dorm-detail',
  templateUrl: './dorm-detail.html',
  standalone: true,
  imports:[CommonModule,FormsModule]
})
export class DormDetailComponent implements OnInit {
  dorm?: Dorm;

  constructor(
    private route: ActivatedRoute,
    private dormsService: DormService
  ) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.dormsService.getDormById(id).subscribe((data) => {
        this.dorm = data;
      });
    }
  }
}
