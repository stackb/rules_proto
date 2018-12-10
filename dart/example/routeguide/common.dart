// Copyright (c) 2017, the gRPC project authors. Please see the AUTHORS file
// for details. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import 'dart:convert';
import 'dart:io';

//
// This path boils down to several things: (1) the bazel package name
// 'routeguide_pb', which is actually the implicit target name of the
// `dart_grpc_compile` rule; (3) the package path of the actual proto
// 'example/proto'; (4) the generated name of the proto library
// 'routeguide.pb.dart'.  
//
import 'routeguide_pb/example/proto/routeguide.pb.dart';

const coordFactor = 1e7;

final List<Feature> featuresDb = _readDatabase();

List<Feature> _readDatabase() {
  final dbData = new File('example/proto/routeguide_features.json').readAsStringSync();
  final List db = jsonDecode(dbData);
  return db.map((entry) {
    final location = new Point()
      ..latitude = entry['location']['latitude']
      ..longitude = entry['location']['longitude'];
    return new Feature()
      ..name = entry['name']
      ..location = location;
  }).toList();
}