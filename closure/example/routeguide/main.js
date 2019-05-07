goog.provide('my.jspb.test.main');

goog.require('jspb.Map');

/**
 * Main entry point for the application.
 * @export
 */
my.jspb.test.main = function() {

  // Increase stacktrace limit in chrome
  Error['stackTraceLimit'] = 150;

  const protoMap = new jspb.Map([]);
  console.log("ok", protoMap);
};
