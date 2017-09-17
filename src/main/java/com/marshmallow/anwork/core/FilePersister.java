package com.marshmallow.anwork.core;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

/**
 * This object loads and saves an object to a file on disk.
 *
 * <p>
 * Created Sep 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class FilePersister<T extends Serializable<?>> implements Persister<T> {

  private final File root;

  /**
   * Initialize a file persister to store stuff at some root.
   *
   * @param root The root at which to persister stuff.
   */
  public FilePersister(File root) {
    this.root = root;
  }

  @Override
  public boolean exists(String context) {
    return convertContextToFile(context).exists();
  }

  @Override
  public void clear(String context) throws IOException {
    if (!exists(context)) {
      throw new IOException("Context does not exist: " + context);
    } else {
      convertContextToFile(context).delete();
    }
  }

  @Override
  public Collection<T> load(String context, Serializer<T> serializer) throws IOException {
    File file = convertContextToFile(context);
    if (!file.exists()) {
      throw new IOException("File for context " + context + " does not exist");
    } else if (!file.canRead()) {
      throw new IOException("Cannot read file for context " + context);
    }

    List<T> ts = new ArrayList<T>();
    try (InputStream inputStream = new FileInputStream(file)) {
      while (inputStream.available() > 0) { // TODO: is this too time intensive?
        T t = serializer.unserialize(inputStream);
        ts.add(t);
      }
    }

    return ts;
  }

  @Override
  public void save(String context, Serializer<T> serializer, Collection<T> data)
      throws IOException {
    File file = convertContextToFile(context);
    if (!file.exists()) {
      file.delete();
    }
    file.getParentFile().mkdirs();
    file.createNewFile();

    try (OutputStream outputStream = new FileOutputStream(file)) {
      for (T t : data) {
        serializer.serialize(t, outputStream);
      }
    }
  }

  private File convertContextToFile(String context) {
    return new File(root, context);
  }
}
