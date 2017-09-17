package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.core.FilePersister;

import java.io.File;

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

  private static final File TEST_RESOURCE_ROOT
      = new File(TestUtilities.TEST_RESOURCES_ROOT, "file-persister-test");

  /**
   * Instantiate this class as a {@link PersisterTest}.
   */
  public FilePersisterTest() {
    super(new FilePersister<Student>(TEST_RESOURCE_ROOT), Student.SERIALIZER);
  }
}
