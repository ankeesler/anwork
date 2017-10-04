package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.core.FilePersister;

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
}
