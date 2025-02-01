import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, throwError } from 'rxjs';
import { collectionURL, httpHeaders } from './common';

@Injectable({
  providedIn: 'root'
})
export class CollectionSearchService {


  searchCollection(collection: string, query: string): Observable<CollectionSearchResponse> {
    let url = collectionURL+collection+"?query="+query
    console.log("Rest to "+url);
    return this.http.get<CollectionSearchResponse>(url, { headers: httpHeaders }).pipe(
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      console.log(error.error.message)

    } else {
      console.log(error.status)
    }
    return throwError(
      console.log('Something is wrong!'));
  };
  constructor(private http: HttpClient) { }
}

export interface CollectionSearchResponse {
  Title: string
  Collection: string
  Query: string
  Documents: Document[]
}

export interface Document {
  Content: string
  Summary: string
  Modified: string
  URL: string
  Title: string
}

/*
{
  "Title": "Search: intranet-ITSKB",
  "Baseurl": "",
  "Version": "0.1.0 (development)",
  "StatusMessage": "Duration 0s - ",
  "Collection": "intranet-ITSKB",
  "Query": "resr",
  "Documents": [
    {
      "Content": "",
      "Modified": "2025-01-09 09:00:05.297 +0100 CET",
      "URL": "https://intranet.unibas.ch/pages/viewpage.action?pageId=311756866",
      "Title": "[Published] M365: RemoteDesktopEnabler  (Article#10871)",
      "IDField": "URL"
    },
    {
      "Content": "",
      "Modified": "2023-12-13 20:01:06.851 +0100 CET",
      "URL": "https://intranet.unibas.ch/pages/viewpage.action?pageId=174744280",
      "Title": "[Published] TeamViewer: Code of Conduct (Article#10431)",
      "IDField": "URL"
    }
  ]
}
*/