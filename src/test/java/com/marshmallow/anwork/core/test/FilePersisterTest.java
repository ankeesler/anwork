package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.core.FilePersister;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import org.junit.After;
import org.junit.Before;
import org.junit.Test;

/**
 * This is a unit test for a {@link FilePersister}.
 *
 * <p>
 * Created Sep 6, 2017
 * </p>
 *
 * @author Andrew
 */
public class FilePersisterTest extends BasePersisterTest {

  private File unreadableFile;

  /**
   * Create the unreadable file for {@link #loadUnreadableFile()} below.
   *
   * @throws IOException if something goes wrong
   */
  @Before
  public void setupUnreadableFile() throws IOException {
    URI unreadableFileDirectory = TestUtilities.getFile(".", FilePersisterTest.class).toURI();
    Path unreadableFilePath = Files.createTempFile(Paths.get(unreadableFileDirectory),
                                                   "FilePersisterTest",
                                                   "");
    unreadableFile = unreadableFilePath.toFile();
    unreadableFile.setReadable(false);
  }

  /**
   * Delete the unreadable file used in {@link #loadUnreadableFile()} below.
   *
   * @throws IOException if something goes wrong
   */
  @After
  public void deleteUnreadableFile() throws IOException {
    unreadableFile.delete(); // I am very confused why this doesn't delete automatically...
  }

  /**
   * Instantiate this class as a {@link PersisterTest}.
   */
  public FilePersisterTest() {
    super(new FilePersister<Student>(TestUtilities.getFile(".", FilePersisterTest.class)),
          Student.SERIALIZER);
  }

  @Test(expected = IOException.class)
  public void loadUnreadableFile() throws IOException {
    persister.load(unreadableFile.getName(), serializer);
  }
}
