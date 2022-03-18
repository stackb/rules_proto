package example.routeguide.scala

import java.util.logging.Logger

import io.grpc.{Server, ServerBuilder}

import example.routeguide.{RouteGuideGrpc}

class RouteGuideServer(server: Server) {

  val logger: Logger = Logger.getLogger(classOf[RouteGuideServer].getName)

  def start(): Unit = {
    server.start()
    logger.info(s"Server started, listening on ${server.getPort}")
    sys.addShutdownHook {
      // Use stderr here since the logger may has been reset by its JVM shutdown hook.
      System.err.println("*** shutting down gRPC server since JVM is shutting down")
      stop()
      System.err.println("*** server shut down")
    }
    ()
  }

  def stop(): Unit = {
    server.shutdown()
  }

  /**
    * Await termination on the main thread since the grpc library uses daemon threads.
    */
  def blockUntilShutdown(): Unit = {
    server.awaitTermination()
  }
}

object RouteGuideServer extends App {
  val features = RouteGuideUtil.parseFeatures(RouteGuideUtil.defaultFeatureFile)

  val server = new RouteGuideServer(
    ServerBuilder
      .forPort(50056)
      .addService(
        RouteGuideGrpc.bindService(
          new RouteGuideService(features),
          scala.concurrent.ExecutionContext.global
        )
      )
      .build()
  )
  server.start()
  server.blockUntilShutdown()
}
