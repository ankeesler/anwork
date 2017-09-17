package com.marshmallow.anwork.core;

import java.io.IOException;

/**
 * This is an object that can be serialized via some medium of type T.
 *
 * <p>
 * This medium may be XML, Google Protocol Buffers, a {@link String}, a custom implementation, etc.
 * </p>
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Serializable<T> {

  /**
   * Turn this object into bytes on an instance of the medium of type T. The instance medium is
   * returned.
   *
   * @return The instance of the medium of type T containing a serialized version of this object;
   *     this should <b>NEVER</b> return <code>null</code>
   * @throws IOException if something goes wrong
   */
  public T marshall() throws IOException;

  /**
   * Read in this object from an instance of the medium of type T.
   *
   * @param t The instance of the medium of type T from which to read this object
   * @throws IOException if something goes wrong
   */
  public void unmarshall(T t) throws IOException;
}
