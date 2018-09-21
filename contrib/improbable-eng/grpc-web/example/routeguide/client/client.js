goog.module('example.routeguide.grpc_web.Client');

const ClientReadableStream = goog.require('grpc.web.ClientReadableStream');
const Feature = goog.require('proto.routeguide.Feature');
const Rectangle = goog.require('proto.routeguide.Rectangle');
const RouteGuideClient = goog.require('proto.routeguide.RouteGuideClient');

/**
 * A skeleton client. The point of this exercise is not actually to create a
 * routeguide client, but show how we can use protobufs in closure code.
 */
class Client {

    /**
     * @param {string} hostname
     * @param {?Object} credentials
     * @param {?Object} options
     */
    constructor(hostname, credentials, options) {
        /**
         * @const @private @type {!RouteGuideClient}
         */
        this.client_ = new RouteGuideClient(hostname, credentials, options);
    }

    /**
     * List existing features.  Promise resolves to a list of features.
     * 
     * @param {!Rectangle} rect
     * @param {!Object} metadata
     * @return {!ClientReadableStream<!Feature>}  
     */
    listFeatures(rect, metadata) {
        return this.client_.listFeatures(rect, metadata);
    }

    /**
     * Run the client.  A real implementation might actually do something here.
     */
    run() {
        console.log("Running grpc-web routeguide client...");
    }

}

exports = Client;