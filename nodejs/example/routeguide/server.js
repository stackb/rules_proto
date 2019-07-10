/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

const featureDb = require('build_stack_rules_proto/example/proto/routeguide_features.json');
const messages = require('build_stack_rules_proto/nodejs/example/routeguide/routeguide/example/proto/routeguide_pb.js')
const services = require('build_stack_rules_proto/nodejs/example/routeguide/routeguide/example/proto/routeguide_grpc_pb.js')

const fs = require('fs');
const path = require('path');
const grpc = require('grpc');
const COORD_FACTOR = 1e7;

/**
 * For simplicity, a point is a record type that looks like
 * {latitude: number, longitude: number}, and a feature is a record type that
 * looks like {name: string, location: point}. feature objects with name===''
 * are points with no feature.
 */

/**
 * List of feature objects at points that have been requested so far.
 */
let feature_list = [];

/**
 * Get a feature object at the given point, or creates one if it does not exist.
 * @param {point} point The point to check
 * @return {feature} The feature object at the point. Note that an empty name
 *     indicates no feature
 */
function checkFeature(point) {
  let feature;
  // Check if there is already a feature object for the given point
  for (let i = 0; i < feature_list.length; i++) {
    feature = feature_list[i];
    if (feature.getLocation().getLatitude() === point.getLatitude() &&
        feature.getLocation().getLongitude() === point.getLongitude()) {
      return feature;
    }
  }
  const name = '';
  feature = new messages.Feature();
  feature.setName(name);
  feature.setLocation(point);
  return feature;
}

/**
 * getFeature request handler. Gets a request with a point, and responds with a
 * feature object indicating whether there is a feature at that point.
 * @param {EventEmitter} call Call object for the handler to process
 * @param {function(Error, feature)} callback Response callback
 */
function getFeature(call, callback) {
  callback(null, checkFeature(call.request));
}

/**
 * listFeatures request handler. Gets a request with two points, and responds
 * with a stream of all features in the bounding box defined by those points.
 * @param {Writable} call Writable stream for responses with an additional
 *     request property for the request value.
 */
function listFeatures(call) {
  const lo = call.request.getLo();
  const hi = call.request.getHi();
  const left = Math.min([lo.getLongitude(), hi.getLongitude()]);
  const right = Math.max([lo.getLongitude(), hi.getLongitude()]);
  const top = Math.max([lo.getLatitude(), hi.getLatitude()]);
  const bottom = Math.min([lo.getLatitude(), hi.getLatitude()]);
  // For each feature, check if it is in the given bounding box
  feature_list.forEach(function(feature) {
    if (feature.getName() === '') {
      return;
    }
    if (feature.getLocation().getLongitude() >= left &&
        feature.getLocation().getLongitude() <= right &&
        feature.getLocation().getLatitude() >= bottom &&
        feature.getLocation().getLatitude() <= top) {
      call.write(feature);
    }
  });
  call.end();
}

/**
 * Calculate the distance between two points using the "haversine" formula.
 * The formula is based on http://mathforum.org/library/drmath/view/51879.html.
 * @param start The starting point
 * @param end The end point
 * @return The distance between the points in meters
 */
function getDistance(start, end) {
  function toRadians(num) {
    return num * Math.PI / 180;
  }
  const R = 6371000;  // earth radius in metres
  const lat1 = toRadians(start.getLatitude() / COORD_FACTOR);
  const lat2 = toRadians(end.getLatitude() / COORD_FACTOR);
  const lon1 = toRadians(start.getLongitude() / COORD_FACTOR);
  const lon2 = toRadians(end.getLongitude() / COORD_FACTOR);

  const deltalat = lat2-lat1;
  const deltalon = lon2-lon1;
  const a = Math.sin(deltalat/2) * Math.sin(deltalat/2) +
      Math.cos(lat1) * Math.cos(lat2) *
      Math.sin(deltalon/2) * Math.sin(deltalon/2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
  return R * c;
}

/**
 * recordRoute handler. Gets a stream of points, and responds with statistics
 * about the "trip": number of points, number of known features visited, total
 * distance traveled, and total time spent.
 * @param {Readable} call The request point stream.
 * @param {function(Error, routeSummary)} callback The callback to pass the
 *     response to
 */
function recordRoute(call, callback) {
  let pointCount = 0;
  let featureCount = 0;
  let distance = 0;
  let previous = null;
  // Start a timer
  const start_time = process.hrtime();
  call.on('data', function(point) {
    pointCount += 1;
    if (checkFeature(point).name !== '') {
      featureCount += 1;
    }
    /* For each point after the first, add the incremental distance from the
     * previous point to the total distance value */
    if (previous != null) {
      distance += getDistance(previous, point);
    }
    previous = point;
  });
  call.on('end', function() {
    const summary = new messages.RouteSummary();
    summary.setPointCount(pointCount);
    summary.setFeatureCount(featureCount);
    // Cast the distance to an integer
    summary.setDistance(distance|0);
    // End the timer
    summary.setElapsedTime(process.hrtime(start_time)[0]);
    callback(null, summary);
  });
}

const route_notes = {};

/**
 * Turn the point into a dictionary key.
 * @param {point} point The point to use
 * @return {string} The key for an object
 */
function pointKey(point) {
  return point.getLatitude() + ' ' + point.getLongitude();
}

/**
 * routeChat handler. Receives a stream of message/location pairs, and responds
 * with a stream of all previous messages at each of those locations.
 * @param {Duplex} call The stream for incoming and outgoing messages
 */
function routeChat(call) {
  call.on('data', function(note) {
    const key = pointKey(note.getLocation());
    /* For each note sent, respond with all previous notes that correspond to
     * the same point */
    if (route_notes.hasOwnProperty(key)) {
      route_notes[key].forEach(function(note) {
        call.write(note);
      });
    } else {
      route_notes[key] = [];
    }
    // Then add the new note to the list
    route_notes[key].push(note);
  });
  call.on('end', function() {
    call.end();
  });
}

/**
 * Get a new server with the handler functions in this file bound to the methods
 * it serves.
 * @return {Server} The new server object
 */
function getServer() {
  const server = new grpc.Server();
  server.addService(services.RouteGuideService, {
    getFeature: getFeature,
    listFeatures: listFeatures,
    recordRoute: recordRoute,
    routeChat: routeChat
  });
  return server;
}

if (require.main === module) {
  let port = '50051';
  if (process.env.SERVER_PORT) {
    port = process.env.SERVER_PORT;
  }
  const addr = '0.0.0.0:'+port;
  // If this is run as a script, start a server on an unused port
  const routeServer = getServer();
  routeServer.bind(addr, grpc.ServerCredentials.createInsecure());

  // Transform the loaded features to Feature objects
  feature_list = featureDb.map(function(value) {
    const feature = new messages.Feature();
    feature.setName(value.name);
    const location = new messages.Point();
    location.setLatitude(value.location.latitude);
    location.setLongitude(value.location.longitude);
    feature.setLocation(location);
    return feature;
  });

  console.log(`Feature database contains ${feature_list.length} entries.`);
  console.log(`Node server listening at ${addr}...`)
  routeServer.start();
}

exports.getServer = getServer;
