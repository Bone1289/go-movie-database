import {Component, OnInit} from '@angular/core';
import {Router} from "@angular/router";

import {Location} from "@angular/common";
import {HttpClient} from "@angular/common/http";

@Component({
  selector: 'app-movies',
  templateUrl: './movies.component.html',
  styleUrls: ['./movies.component.css']
})
export class MoviesComponent implements OnInit {

  public movies: any;
  private baseURL: string = "http://localhost:12345/movies";

  constructor(private http: HttpClient,
              private router: Router,
              private location: Location) {
  }

  ngOnInit(): void {
    this.location.subscribe(() => {
      this.refresh();
    });
    this.refresh();
  }

  private refresh() {
    this.http.get(this.baseURL)
      .subscribe(result => {
        console.log("Testing")
        this.movies = result;
      })
  }

  public search(event: any) {
    let url = this.baseURL;
    if (event.target.value) {
      url = this.baseURL + event.target.value;
    }
    this.http.get(url).subscribe(result => {
      console.log("Testing")
      this.movies = result;
    })
  }

  public create() {
    this.router.navigate(["create"]);
  }

}
