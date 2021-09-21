goog.module('example.routeguide.closure.GrpcJsClient');

const Feature = goog.require('proto.example.routeguide.Feature');
const GoogPromise = goog.require('goog.Promise');
const Rectangle = goog.require('proto.example.routeguide.Rectangle');
const RouteguideClient = goog.require('proto.example.routeguide.RouteguideClient');

/**
 * A skeleton client. The point of this exercise is not actually to create a
 * routeguide client, but show how we can use protobufs in closure code.
 */
class GrpcJsClient {
    constructor() {
        /**
         * @const @private @type {!RouteguideClient}
         */
        this.client_ = new RouteguideClient();
    }

    /**
     * List existing features.  Promise resolves to a list of features.
     *
     * @param {!Rectangle} rect
     * @return {!GoogPromise<!Array<!Feature>>} feature
     */
    listFeatures(rect) {
        /**
         * @type {!Array<!Feature>}
         */
        const features = [];

        return this.client_.getRouteGuide()
            .listFeatures(rect, f => features.push(f))
            .then(() => features);
    }

    /**
     * Run the client.  A real implementation might actually do something here.
     */
    run() {
        console.log("Running grpc.js routeguide client...");
    }
}
exports = GrpcJsClient;