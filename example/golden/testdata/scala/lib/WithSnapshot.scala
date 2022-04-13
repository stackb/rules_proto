package lib

import proto.Complete

/** 
  * An interface intended for extending ScalaPB generated case classes, but may
  * be used more generically.
  *
  * @tparam Message The type used for either snapshot messages or updates.
  * @tparam A       The class to which this interface is being added.
  */
trait WithSnapshot[Message, A] {
  def withMessage(message: Message): A

  def withComplete(complete: Complete): A

  def getMessage: Message

  def getComplete: Complete

  final def toEither(isComplete: PartialFunction[WithSnapshot[_, _], Boolean]): Either[Message, Complete] =
    if (isComplete.lift(this).contains(true)) {
      Right(getComplete)
    } else {
      Left(getMessage)
    }
}
