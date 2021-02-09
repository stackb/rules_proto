goog.module('example.routeguide.closure');

const Point = goog.require('proto.routeguide.Point');
// goog.require('jspb.Map');

class Application {
  constructor() {
  }

  sayHello() {
    const point = new Point();
    point.setLatitude(413069058);
    point.setLongitude(-744597778);

    console.log('The point is', point);
  }
}

/**
  * Entrypoint function that is exported.
 */
function main() {
  const app = new Application();
  app.sayHello();
}
