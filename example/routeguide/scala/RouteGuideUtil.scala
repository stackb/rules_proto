package example.routeguide.scala

import java.net.URL
import java.util.logging.Logger
import scalapb.json4s.JsonFormat
import scala.io.Source

import example.routeguide.{Feature, FeatureDatabase, Point}

object RouteGuideUtil {
  val logger: Logger = Logger.getLogger(getClass.getName)
  val defaultFeatureFile: URL = getClass.getClassLoader.getResource("example/routeguide/feature_db.json")

  /**
    * Get a canned sequence of features so we don't have to parse the json file.
    */
  def getFeatures(): Seq[Feature] = {
    val features: Seq[Feature] = Seq(
      Feature(
        name = "Patriots Path, Mendham, NJ 07945, USA",
        location = Some(Point(407838351, -746143763))),
      Feature(
        name = "101 New Jersey 10, Whippany, NJ 07981, USA",
        location = Some(Point(408122808, -743999179)))
    )
    features
  }

  /**
    * Parses the JSON input file containing the list of features.
    */
  def parseFeatures(file: URL): Seq[Feature] = {
    logger.info(s"Loading features from ${file}")
    var features: Seq[Feature] = Seq.empty
    val input = file.openStream
    try {
      val source = Source.fromInputStream(input)
      try {
        val db = JsonFormat.fromJsonString[FeatureDatabase](source.getLines().mkString("\n"))
        features = db.feature
        logger.info(s"Parsed features from ${file.getPath}")
      } finally source.close()
    } finally input.close
    logger.info(s"Loaded ${features.size} features")
    features
  }

}

object RouteGuideServiceUtil {
  def isValid(feature: Feature): Boolean = feature.name.nonEmpty

  val COORD_FACTOR: Double = 1e7
  def getLatitude(point: Point): Double = point.latitude.toDouble / COORD_FACTOR
  def getLongitude(point: Point): Double = point.longitude.toDouble / COORD_FACTOR
  /**
    * Calculate the distance between two points using the "haversine" formula.
    * This code was taken from http://www.movable-type.co.uk/scripts/latlong.html.
    *
    * @param start The starting point
    * @param end   The end point
    * @return The distance between the points in meters
    */
  def calcDistance(start: Point, end: Point) = {
    val lat1 = getLatitude(start)
    val lat2 = getLatitude(end)
    val lon1 = getLongitude(start)
    val lon2 = getLongitude(end)
    val r = 6371000
    // meters
    import Math._
    val phi1 = toRadians(lat1)
    val phi2 = toRadians(lat2)
    val deltaPhi = toRadians(lat2 - lat1)
    val deltaLambda = toRadians(lon2 - lon1)
    val a = Math.sin(deltaPhi / 2) * Math.sin(deltaPhi / 2) + cos(phi1) * cos(phi2) * sin(deltaLambda / 2) * sin(deltaLambda / 2)
    val c = 2 * atan2(sqrt(a), sqrt(1 - a))
    (r * c).toInt
  }
}