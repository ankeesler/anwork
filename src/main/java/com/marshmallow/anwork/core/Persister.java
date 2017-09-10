package com.marshmallow.anwork.core;

import java.io.IOException;
import java.util.Collection;

/**
 * This represents an object that can persist some data.
 *
 * This object operates on some loosely typed "context" object. This can be a
 * file path, a URI, simply a file name, or anything else. The context is
 * loosely typed on purpose for academic purposes.
 *
 * @author Andrew
 * Created Sep 4, 2017
 * @see Serializer
 */
public interface Persister<T> {

  /**
   * Returns whether or not the context exists.
   *
   * @param context The context to check
   * @return Whether or not the context exists
   * @throws IOException If something goes wrong
   */
  public boolean contextExists(String context) throws IOException;

  /**
   * Load some objects from a persistent store.
   *
   * @param context The context from which to load the objects; if this context
   * does not exist, then the persister should throw an exception
   * @param serializer The serializer to use when loading the objects
   * @return An array of objects
   * @throws IOException if something goes wrong, like if the
   * context does not exist
   * @see #contextExists(String)
   */
  public Collection<T> load(String context, Serializer<T> serializer) throws IOException;

  /**
   * Save some objects to a persistent store.
   *
   * @param context The context to save the objects to
   * @param serializer The serializer to use
   * @param data The objects to save
   * @throws IOException if something goes wrong
   */
  public void save(String context, Serializer<T> serializer, Collection<T> data) throws IOException;
}
