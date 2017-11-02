package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.core.FilePersister;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

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

  /**
   * Instantiate this class as a {@link PersisterTest}.
   */
  public FilePersisterTest() {
    super(new FilePersister<Student>(TestUtilities.getFile(".", FilePersisterTest.class)),
          Student.SERIALIZER);
  }

  @Test(expected = IOException.class)
  public void loadUnreadableFile() throws IOException {
    URI unreadableFileDirectory = TestUtilities.getFile(".", FilePersisterTest.class).toURI();
    Path unreadableFilePath = Files.createTempFile(Paths.get(unreadableFileDirectory),
                                                   "FilePersisterTest",
                                                   "");
    File unreadableFile = unreadableFilePath.toFile();
    unreadableFile.setReadable(false);
    persister.load(unreadableFile.getName(), serializer);
  }
}
