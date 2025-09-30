import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
import { Dorm } from './dorm';

export interface DormPrice {
  name: string;
  price: number;
}
export interface FreeSpots {
  city: string;
  freeSpots: number;
}

export interface CityDorms {
  city: string;
  dorms: DormPrice[];
  averagePrice: number;
}

export interface FacilitiesSummary {
  [city: string]: {
    [facility: string]: number;
  };
}

// dorm.ts ili open-dorm.ts
export interface OpenDorm {
  id: string;
  name: string;
  city: string;
  capacity: number;
  occupied: number;
  occupancy_rate: number;
  type: string;
  amenities: string[];
  average_rating: number;
  comments_count: number;
  tags: string[];
  price: number;
}



@Injectable({
  providedIn: 'root'
})
export class Opendata {

  private baseUrl = 'http://localhost:8000/opendata';

  constructor(private http: HttpClient) {}
  
  getDormsWithAvgPrice(): Observable<{ [city: string]: number }> {
  return this.http.get<{ [city: string]: number }>(`${this.baseUrl}/avg-price`);
}

  getDorms(): Observable<Dorm[]> {
  return this.http.get<Dorm[]>(`${this.baseUrl}/open-dorms`);
}
  getAvgRating(): Observable<OpenDorm[]> {
  return this.http.get<OpenDorm[]>(`${this.baseUrl}/open-dorms`);
}

 getFreeSpots(): Observable<FreeSpots[]> {
    return this.http.get<{ [city: string]: number }>(`${this.baseUrl}/cities/free-spots`).pipe(
      map(data =>
        Object.entries(data).map(([city, freeSpots]) => ({ city, freeSpots }))
      )
    );
  }

getOccupancyPerCity(): Observable<Record<string, number>> {
  return this.http.get<Record<string, number>>(`${this.baseUrl}/trends/occupancy`);
}

getFacilitiesSummary(): Observable<FacilitiesSummary> {
    return this.http.get<FacilitiesSummary>(`${this.baseUrl}/facilities/summary`);
  }

  
}
