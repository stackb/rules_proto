package example

import java.net.URL
import java.util.logging.Logger

import scalapb.json4s.JsonFormat

import scala.io.Source

import io.grpc.examples.routeguide.routeguide.{Feature, FeatureDatabase, Point}

object RouteGuidePersistence {
  val logger: Logger = Logger.getLogger(getClass.getName)

  // val defaultFeatureFile: URL = getClass.getClassLoader.getResource("/scala/example/routeguide/route_guide_db.json")
  val defaultFeatureFile: URL = getClass.getClassLoader.getResource("scala/example/routeguide/route_guide_db.json")

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
