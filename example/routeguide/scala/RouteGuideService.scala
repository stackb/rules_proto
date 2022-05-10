package example.routeguide.scala

import java.util.logging.Logger
import java.util.concurrent.TimeUnit.NANOSECONDS

import io.grpc.stub.StreamObserver
import scala.concurrent.Future

import example.routeguide.{Feature, Point, RouteNote, RouteSummary, Rectangle, RouteGuideGrpc}

class RouteGuideService(features: Seq[Feature]) extends RouteGuideGrpc.RouteGuide {

  val logger: Logger = Logger.getLogger(classOf[RouteGuideService].getName)

  private val routeNotes: AtomicRef[Map[Point, Seq[RouteNote]]] = new AtomicRef(Map.empty)

  /**
    * Gets the {@link Feature} at the requested {@link Point}. If no feature at
    * that location exists, an unnamed feature is returned at the provided
    * location.
    *
    * @param request the requested location for the feature.
    */
  override def getFeature(request: Point): Future[Feature] = {
    Future.successful(findFeature(request))
  }

  /**
    * Gets all features contained within the given bounding {@link Rectangle}.
    *
    * @param request          the bounding rectangle for the requested features.
    * @param responseObserver the observer that will receive the features.
    */
  override def listFeatures(request: Rectangle, responseObserver: StreamObserver[Feature]): Unit = {
    val left = Math.min(request.getLo.longitude, request.getHi.longitude)
    val right = Math.max(request.getLo.longitude, request.getHi.longitude)
    val top = Math.max(request.getLo.latitude, request.getHi.latitude)
    val bottom = Math.min(request.getLo.latitude, request.getHi.latitude)

    features.foreach { feature =>
      if (RouteGuideServiceUtil.isValid(feature)) {
        val lat = feature.getLocation.latitude
        val lon = feature.getLocation.longitude
        if (lon >= left && lon <= right && lat >= bottom && lat <= top) {
          responseObserver.onNext(feature)
        }
      }
    }
    responseObserver.onCompleted()
  }

  /**
    * Gets a stream of points, and responds with statistics about the "trip":
    * number of points, number of known features visited, total distance
    * traveled, and total time spent.
    *
    * @param responseObserver an observer to receive the response summary.
    * @return an observer to receive the requested route points.
    */
  override def recordRoute(responseObserver: StreamObserver[RouteSummary]): StreamObserver[Point] =
    new StreamObserver[Point] {
      var pointCount: Int = 0
      var featureCount: Int = 0
      var distance: Int = 0
      var previous: Option[Point] = None
      var startTime: Long = System.nanoTime

      override def onNext(point: Point): Unit = {
        pointCount += 1
        if (RouteGuideServiceUtil.isValid(findFeature(point))) {
          featureCount += 1
        }
        // For each point after the first, add the incremental distance from the
        // previous point to the total distance value.
        previous.foreach{ prev =>
          distance += RouteGuideServiceUtil.calcDistance(prev, point)
        }
        previous = Some(point)
      }
      override def onCompleted(): Unit = {
        val seconds = NANOSECONDS.toSeconds(System.nanoTime - startTime)
        responseObserver.onNext(RouteSummary(pointCount, featureCount, distance, seconds.toInt))
        responseObserver.onCompleted
      }

      override def onError(t: Throwable): Unit =
        logger.warning("recordRoute cancelled")
    }

  /**
    * Receives a stream of message/location pairs, and responds with a stream of
    * all previous messages at each of those locations.
    *
    * @param responseObserver an observer to receive the stream of previous
    * messages.
    * @return an observer to handle requested message/location pairs.
    */
  override def routeChat(responseObserver: StreamObserver[RouteNote]): StreamObserver[RouteNote] =
    new StreamObserver[RouteNote]() {
      override def onNext(note: RouteNote): Unit = {
        val notes = getNotes(note.getLocation)
        // Respond with all previous notes at this location
        notes.foreach(responseObserver.onNext)
        // Now add the new note to the list
        addNote(note)
      }

      override def onError(t: Throwable): Unit = {
        logger.warning("routeChat cancelled")
      }

      override def onCompleted(): Unit = {
        responseObserver.onCompleted
      }
    }

  private def findFeature(point: Point): Feature = {
    features.find { feature =>
      feature.getLocation.latitude == point.latitude && feature.getLocation.longitude == point.longitude
    } getOrElse new Feature(location = Some(point))
  }

  private def getNotes(point: Point): Seq[RouteNote] = {
    routeNotes.get.getOrElse(point, Seq.empty)
  }

  private def addNote(note: RouteNote): Unit = {
    routeNotes.updateAndGet { notes =>
      val existingNotes = notes.getOrElse(note.getLocation, Seq.empty)
      val updatedNotes = existingNotes :+ note
      notes + (note.getLocation -> updatedNotes)
    }
  }
}
