package com.marshmallow.anwork.core;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;

/**
 * This is an object that is able to take an instance of a particular type and turn it into an array
 * of bytes.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Serializer<T extends Serializable<?>> {

  /**
   * Turn the provided object into an array of bytes and add those bytes to an {@link OutputStream}.
   *
   * @param t The provided object to turn into bytes
   * @param outputStream The {@link OutputStream} where the bytes should be written
   * @throws IOException if something goes wrong
   */
  public void serialize(T t, OutputStream outputStream) throws IOException;

  /**
   * Read one object from an array of bytes in the form of an {@link InputStream}.
   *
   * @param inputStream The {@link InputStream} from where to read the bytes
   * @return An instance of the object to read from the stream; this is never <code>null</code>!
   * @throws IOException if something goes wrong, like there is no object to be read.
   */
  public T unserialize(InputStream inputStream) throws IOException;
}
