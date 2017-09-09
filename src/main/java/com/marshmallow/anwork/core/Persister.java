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
 * @date Sep 4, 2017
 * @see Serializer
 */
public interface Persister<T> {

  /**
   * Load some objects from a persistent store.
   *
   * @param context The context from which to load the objects
   * @param serializer The serializer to use when loading the objects
   * @return An array of objects
   * @throws An {@link IOException} if something goes wrong
   */
  public Collection<T> load(String context, Serializer<T> serializer) throws IOException;

  /**
   * Save some objects to a persistent store.
   *
   * @param context The context to save the objects to
   * @param serializer The serializer to use
   * @param data The objects to save
   * @throws An {@link IOException} if something goes wrong
   */
  public void save(String context, Serializer<T> serializer, Collection<T> data) throws IOException;
}