package com.marshmallow.anwork.core;

import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.io.LineNumberReader;
import java.io.Writer;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

/**
 * This object loads and saves an object to a file on disk.
 *
 * @author Andrew
 * @date Sep 4, 2017
 */
public class FilePersister<T> implements Persister<T> {

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
  public Collection<T> load(String context, Serializer<T> serializer) throws IOException {
    File file = convertContextToFile(context);
    if (!file.exists()) {
      throw new IOException("File for context " + context + " does not exist");
    } else if (!file.canRead()) {
      throw new IOException("Cannot read file for context " + context);
    }

    List<T> ts = new ArrayList<T>();
    try (LineNumberReader lineNumberReader = new LineNumberReader(new FileReader(file))) {
      String line;
      while ((line = lineNumberReader.readLine()) != null) {
        T t = serializer.unmarshall(line);
        if (t == null) {
          throw new IOException("Unknown serialization '" + line + "' for serializer " + serializer);
        }
        ts.add(t);
      }
    } catch (IOException ioe) {
      throw ioe;
    }

    return ts;
  }

  @Override
  public void save(String context, Serializer<T> serializer, Collection<T> data) throws IOException {
    File file = convertContextToFile(context);
    if (!file.exists()) {
      file.delete();
    }
    file.getParentFile().mkdirs();
    file.createNewFile();

    try (Writer fileWriter = new FileWriter(file)){
      for (T t : data) {
        String marshalled = serializer.marshall(t);
        if (marshalled == null) {
          throw new IOException("Cannot serialize '" + t + "' with serializer " + serializer);
        }
        fileWriter.append(marshalled);
        fileWriter.append("\n");
      }
    } catch (IOException ioe) {
      throw ioe;
    }
  }

  private File convertContextToFile(String context) {
    return new File(root, context);
  }
}
