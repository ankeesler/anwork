package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;

import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.core.test.BaseSerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskState;
import com.marshmallow.anwork.task.protobuf.TaskStateProtobuf;

import org.junit.Test;

/**
 * A {@link BaseSerializerTest} for {@link Task} objects.
 *
 * <p>
 * Created Sep 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskSerializerTest extends BaseSerializerTest<Task> {

  /**
   * Instantiate this test as a {@link BaseSerializerTest}.
   */
  protected Serializer<Task> getSerializer() {
    return Task.SERIALIZER;
  }

  @Test
  public void testThatTaskStateEnumsLineUp() {
    // This test is to ensure that our TaskProtobuf.TaskState enum is the same
    // as our TaskState enum.
    assertEquals(TaskState.WAITING.ordinal(), TaskStateProtobuf.WAITING.ordinal());
    assertEquals(TaskState.BLOCKED.ordinal(), TaskStateProtobuf.BLOCKED.ordinal());
    assertEquals(TaskState.RUNNING.ordinal(), TaskStateProtobuf.RUNNING.ordinal());
    assertEquals(TaskState.FINISHED.ordinal(), TaskStateProtobuf.FINISHED.ordinal());
  }
}
