goog.provide('my.jspb.test.main');

goog.require('jspb.Map');

/**
 * Main entry point for the application. All we basically checking is that the
 * closure compiler is satisfied and can compile the bundle.
 *
 * @export
 */
my.jspb.test.main = function() {

  // Increase stacktrace limit in chrome
  Error['stackTraceLimit'] = 150;

  const protoMap = new jspb.Map([]);
  console.log("ok", protoMap);
};
