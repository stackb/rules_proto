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

const async = require('async');
const _ = require('lodash');
const grpc = require('grpc');

const messages = require('build_stack_rules_proto/nodejs/example/routeguide/routeguide/example/proto/routeguide_pb.js')
const services = require('build_stack_rules_proto/nodejs/example/routeguide/routeguide/example/proto/routeguide_grpc_pb.js')

// This is included as data in the client, so we can load
// this database as a constant.
const featureList = require('build_stack_rules_proto/example/proto/routeguide_features.json');
console.log(`Loaded ${featureList.length} from feature database`);

let port = '50051';
if (process.env.SERVER_PORT) {
  port = process.env.SERVER_PORT;
}
const addr = 'localhost:' + port;

const client = new services.RouteGuideClient(
  addr,
  grpc.credentials.createInsecure());

const COORD_FACTOR = 1e7;

/**
 * Run the getFeature demo. Calls getFeature with a point known to have a
 * feature and a point known not to have a feature.
 * @param {function} callback Called when this demo is complete
 */
function runGetFeature(callback) {
  const next = _.after(2, callback);
  function featureCallback(error, feature) {
    if (error) {
      console.warn("ERROR occured while attempting to get feature", error);
      if (callback) {
        callback(error);
        callback = undefined
      }
      return;
    }

    // console.log("Feature", feature.toObject());

    if (feature.getName() && feature.getName() != "undefined") {
      console.log('Found feature called "' + feature.getName() + '" at ' +
          feature.getLocation().getLatitude()/COORD_FACTOR + ', ' +
          feature.getLocation().getLongitude()/COORD_FACTOR);
    } else {
      console.log('Found no feature at ' +
          feature.getLocation().getLatitude()/COORD_FACTOR + ', ' +
          feature.getLocation().getLongitude()/COORD_FACTOR);
    }
    next();
  }
  const point1 = newPoint(409146138, -746188906);
  const point2 = newPoint(1, 1);
  client.getFeature(point1, featureCallback);
  client.getFeature(point2, featureCallback);
}

/**
 * @param {number} latitude
 * @param {number} longitude
 * @returns {!messages.Point}
 */
function newPoint(latitude, longitude) {
  const point = new messages.Point()
  point.setLatitude(latitude);
  point.setLongitude(longitude);
  return point;
}

/**
 * @param {!messages.Point} lo
 * @param {!messages.Point} hi
 * @returns {!messages.Rectangle}
 */
function newRectangle(lo, hi) {
  const rect = new messages.Rectangle()
  rect.setLo(lo);
  rect.setHi(hi);
  return rect;
}

/**
 * @param {!messages.Point} point
 * @param {string} message
 * @returns {!messages.Note}
 */
function newNote(point, message) {
  const note = new messages.RouteNote()
  note.setLocation(point);
  note.setMessage(message);
  return note;
}

/**
 * Run the listFeatures demo. Calls listFeatures with a rectangle containing all
 * of the features in the pre-generated database. Prints each response as it
 * comes in.
 * @param {function} callback Called when this demo is complete
 */
function runListFeatures(callback) {
  const rectangle = newRectangle(
    newPoint(400000000, -750000000),
    newPoint(420000000, -73000000));
  console.log('Looking for features between 40, -75 and 42, -73');
  const call = client.listFeatures(rectangle);
  call.on('data', function(feature) {
    console.log('Found feature called "' + feature.getName() + '" at ' +
          feature.getLocation().getLatitude()/COORD_FACTOR + ', ' +
          feature.getLocation().getLongitude()/COORD_FACTOR);
  });
  call.on('end', callback);
}

/**
 * Run the recordRoute demo. Sends several randomly chosen points from the
 * pre-generated feature database with a constiable delay in between. Prints the
 * statistics when they are sent from the server.
 * @param {function} callback Called when this demo is complete
 */
function runRecordRoute(callback) {

  const num_points = 10;

  const call = client.recordRoute(function(error, stats) {
    if (error) {
      callback(error);
      return;
    }
    console.log('Finished trip with', stats.getPointCount(), 'points');
    console.log('Passed', stats.getFeatureCount(), 'features');
    console.log('Traveled', stats.getDistance(), 'meters');
    console.log('It took', stats.getElapsedTime(), 'seconds');
    callback();
  });

  /**
   * Constructs a function that asynchronously sends the given point and then
   * delays sending its callback
   * @param {number} lat The latitude to send
   * @param {number} lng The longitude to send
   * @return {function(function)} The function that sends the point
   */
  function pointSender(lat, lng) {
    /**
     * Sends the point, then calls the callback after a delay
     * @param {function} callback Called when complete
     */
    return function(callback) {
      console.log('Visiting point ' + lat/COORD_FACTOR + ', ' +
          lng/COORD_FACTOR);
      call.write(newPoint(lat, lng));
      _.delay(callback, _.random(100, 500));
    };
  }

  const pointSenders = [];

  for (let i = 0; i < num_points; i++) {
    const randIndex = _.random(0, featureList.length - 1)
    console.log("randomIndex", randIndex);
    const randomPointJson = featureList[randIndex];
    const randomPoint = newPoint(randomPointJson.location.latitude, randomPointJson.location.longitude)
    console.log("randomPoint", randomPointJson, randomPoint.toObject());
    pointSenders[i] = pointSender(randomPoint.getLatitude(),
                                    randomPoint.getLongitude());
  }

  async.series(pointSenders, function() {
    call.end();
  });

}

/**
 * Run the routeChat demo. Send some chat messages, and print any chat messages
 * that are sent from the server.
 * @param {function} callback Called when the demo is complete
 */
function runRouteChat(callback) {
  const call = client.routeChat();
  call.on('data', function(note) {
    console.log('Got message "' + note.getMessage() + '" at ' +
        note.getLocation().getLatitude() + ', ' + note.getLocation().getLongitude());
  });

  call.on('end', callback);

  const notes = [
    newNote(newPoint(0, 0), 'First message'),
    newNote(newPoint(0, 1), 'Second message'),
    newNote(newPoint(1, 0), 'Third message'),
    newNote(newPoint(0, 0), 'Fourth message'),
  ]
  for (let i = 0; i < notes.length; i++) {
    const note = notes[i];
    console.log('Sending message "' + note.getMessage() + '" at ' +
        note.getLocation().getLatitude() + ', ' + note.getLocation().getLongitude());
    call.write(note);
  }
  call.end();
}

/**
 * Run all of the demos in order
 */
function main() {
  client.waitForReady(4000, () => {
    async.series([
      runGetFeature,
      runListFeatures,
      runRecordRoute,
      runRouteChat
    ]).catch((e) => {
      console.log(e)
      process.exit(1);
    });
  });
}

if (require.main === module) {
  main();
}

exports.runGetFeature = runGetFeature;

exports.runListFeatures = runListFeatures;

exports.runRecordRoute = runRecordRoute;

exports.runRouteChat = runRouteChat;
