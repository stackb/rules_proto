package example

import java.net.URL
import java.util.logging.Logger

import scalapb.json4s.JsonFormat

import scala.io.Source

import io.grpc.examples.routeguide.routeguide.{Feature, FeatureDatabase}

object RouteGuidePersistence {
  val logger: Logger = Logger.getLogger(getClass.getName)

  val defaultFeatureFile: URL = getClass.getClassLoader.getResource("/scala/example/routeguide/route_guide_db.json")

  /**
    * Parses the JSON input file containing the list of features.
    */
  def parseFeatures(file: URL): Seq[Feature] = {
    logger.info(s"Loading features from ${file.getPath}")
    var features: Seq[Feature] = Seq.empty
    val input = file.openStream
    try {
      val source = Source.fromInputStream(input)
      try {
        features = JsonFormat.fromJsonString[FeatureDatabase](source.getLines().mkString("\n")).feature
      } finally source.close()
    } finally input.close
    logger.info(s"Loaded ${features.size} features")
    features
  }

}
