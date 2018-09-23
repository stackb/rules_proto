// Copyright 2018 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

extern crate routeguide;
extern crate grpc;
extern crate tls_api_stub;

use std::thread;
use std::env;
use std::str::FromStr;
use std::iter;
use futures;

use routeguide::*;

struct RouteGuideImpl;

impl RouteGuide for RouteGuideImpl {
  fn list_features(&self, _m: grpc::RequestOptions, req: Rectangle) -> grpc::StreamingResponse<Feature> {
    // Send back 13 dummy list response objects
    let iter = iter::repeat(()).map(|_| {
      let s = "MyTestFileName.bin".to_owned();
      let mut feature = Feature::new();
      feature.set_name(s);
      feature
    }).take(13);
    grpc::StreamingResponse::iter(iter);
  }

  fn record_route(&self, o: grpc::RequestOptions, p: grpc::StreamingRequest<Point>) -> grpc::SingleResponse<RouteSummary> {
    // let result = p.into_iter(() -> {

    // });
    //  {
    //   let summary = RouteSummary::new();
    //   summary
    // };
    // let iter = iter::repeat(()).map(|_| {
    //   let mut summary = RouteSummary::new();
    //   summary
    // }).take(13);

    // match p.iter() {
    //   Err(e) => panic!("{:?}", e),
    //   Ok((_, stream)) => {
    //     for item in stream {
    //       let point = item.unwrap();
    //       println!("> {}", point);
    //     }
    //   }
    // }
    let promise = futures::Future::ok().map(|| {
      let summary = RouteSummary::new();
      summary
    });
    //let future: grpc::GrpcFuture<RouteSummary> = Box::new(promise);

    grpc::SingleResponse::new(promise)
  }

  fn route_chat(&self, o: grpc::RequestOptions, p: grpc::StreamingRequest<RouteNote>) -> grpc::StreamingResponse<RouteNote> {

  }

  fn get_feature(&self, o: grpc::RequestOptions, p: Point) -> grpc::SingleResponse<Feature> {
    let mut r = Feature::new();
    r.set_name(format!("test"));
    grpc::SingleResponse::completed(r)
  }

}

fn main() {
    let mut server = grpc::ServerBuilder::<tls_api_stub::TlsAcceptor>::new();
    let port = u16::from_str(&env::args().nth(1).unwrap_or("50051".to_owned())).unwrap();
    server.http.set_port(port);
    server.add_service(RouteGuideServer::new_service_def(RouteGuideImpl));
    server.http.set_cpu_pool_threads(4);
    let server = server.build().expect("server");
    let port = server.local_addr().port().unwrap();
    println!("RouteGuide server started on port {}", port);

    loop {
        thread::park();
    }
}