import { Component, OnInit } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import { Location } from '@angular/common';

@Component({
  selector: 'app-create',
  templateUrl: './create.component.html',
  styleUrls: ['./create.component.css']
})
export class CreateComponent implements OnInit {

  public movie: any;

  constructor(private location: Location,
              private http: HttpClient) {
    this.movie = {
      "name": "",
      "genre": "",
      "formats": {
        "digital": false,
        "bluray": false,
        "dvd": false
      }
    }
  }

  ngOnInit(): void {
  }

  public save() {
    if(this.movie.name && this.movie.genre) {
      this.http.post("http://localhost:12345/movies", JSON.stringify(this.movie))
        .subscribe(result => {
          this.location.back();
        });
    }
  }
}
