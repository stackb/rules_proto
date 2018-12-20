import com.google.protobuf.compiler.PluginProtos.CodeGeneratorRequest
import com.google.protobuf.ExtensionRegistry
import java.io.{ByteArrayOutputStream, InputStream, PrintWriter, StringWriter}
import scala.reflect.io.Streamable
import scalapb.compiler.ProtobufGenerator
import scalapb.options.compiler.Scalapb

object CompilerPlugin {
  def main(args: Array[String]): Unit = {
    val registry = ExtensionRegistry.newInstance()
    Scalapb.registerAllExtensions(registry)
    val request = CodeGeneratorRequest.parseFrom(Streamable.bytes(System.in), registry)
    System.out.write(ProtobufGenerator.handleCodeGeneratorRequest(request).toByteArray)
  }
}
