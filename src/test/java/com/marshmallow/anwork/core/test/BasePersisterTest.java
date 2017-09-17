package com.marshmallow.anwork.core.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.core.Serializer;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.util.Arrays;
import java.util.Collection;
import java.util.Collections;

import org.junit.Test;

/**
 * This is a generic test for a {@link Persister} implementation.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class BasePersisterTest {

  private static final String DEFAULT_CONTEXT = "default-context";
  private static final String CONTEXT_A = "context-a";
  private static final String CONTEXT_B = "context-b";

  private Persister<Student> persister;
  private Serializer<Student> serializer;

  /**
   * Instatiate this class with a persister and a serializer to test.
   *
   * @param persister The persister to test
   * @param serializer The default serializer to test
   */
  public BasePersisterTest(Persister<Student> persister, Serializer<Student> serializer) {
    this.persister = persister;
    this.serializer = serializer;
  }

  @Test(expected = IOException.class)
  public void testContextDoesNotExist() throws IOException {
    String context = "this-context-does-not-exist";
    assertFalse(persister.exists(context));
    persister.load(context, serializer);
  }

  @Test(expected = IOException.class)
  public void testBadSerialization() throws IOException {
    String context = "this-context-has-a-bad-serialization";
    assertTrue(persister.exists(context));
    persister.load(context, serializer);
  }

  @Test(expected = IOException.class)
  public void testBadClear() throws IOException {
    String context = "this-context-does-not-exist";
    assertFalse(persister.exists(context));
    persister.clear(context);
  }

  @Test(expected = IOException.class)
  public void testBadLoadAfterClear() throws IOException {
    Student saved = new Student("Andrew", 12345);
    persister.save(DEFAULT_CONTEXT, serializer, Collections.singleton(saved));
    persister.clear(DEFAULT_CONTEXT);
    persister.load(DEFAULT_CONTEXT, serializer);
  }

  @Test(expected = IOException.class)
  public void testBadSaveBecauseOfSerializer() throws IOException {
    Student saved = new Student("Andrew", 12345);
    Serializer<Student> badSerializer = new Serializer<Student>() {

      @Override
      public void serialize(Student t, OutputStream outputStream) throws IOException {
        throw new IOException("Failed to serialize!");
      }

      @Override
      public Student unserialize(InputStream inputStream) throws IOException {
        return saved;
      }
    };
    persister.save(DEFAULT_CONTEXT, badSerializer, Collections.singleton(saved));
  }

  @Test(expected = IOException.class)
  public void testBadLoadBecauseOfSerializer() throws IOException {
    Serializer<Student> badSerializer = new Serializer<Student>() {
      @Override
      public void serialize(Student t, OutputStream outputStream) throws IOException {
        // pass
      }

      @Override
      public Student unserialize(InputStream inputStream) throws IOException {
        throw new IOException("Failed to unserialize!");
      }
    };
    assertTrue(persister.exists(CONTEXT_A));
    persister.load(CONTEXT_A, badSerializer);
  }

  @Test()
  public void loadOneTest() throws IOException {
    String context = "load-one-test";
    assertTrue(persister.exists(context));
    Collection<Student> loadeds = persister.load(context, serializer);
    assertEquals(1, loadeds.size());

    Student loaded = loadeds.toArray(new Student[0])[0];
    assertEquals("Tuna Fish Marlin", loaded.getName());
    assertEquals(8675309, loaded.getId());
  }

  @Test()
  public void loadThreeTest() throws IOException {
    String context = "load-three-test";
    assertTrue(persister.exists(context));
    Collection<Student> loadeds = persister.load(context, serializer);
    assertEquals(3, loadeds.size());

    Student[] loadedsArray = loadeds.toArray(new Student[0]);
    assertEquals("Tuna", loadedsArray[0].getName());
    assertEquals(100, loadedsArray[0].getId());
    assertEquals("Fish", loadedsArray[1].getName());
    assertEquals(200, loadedsArray[1].getId());
    assertEquals("Marlin", loadedsArray[2].getName());
    assertEquals(300, loadedsArray[2].getId());
  }

  @Test
  public void saveOneTest() throws IOException {
    Student saved = new Student("Andrew", 12345);
    persister.save(DEFAULT_CONTEXT, serializer, Collections.singleton(saved));

    Collection<Student> loadeds = persister.load(DEFAULT_CONTEXT, serializer);
    assertEquals(1, loadeds.size());

    Student loaded = loadeds.toArray(new Student[0])[0];
    assertEquals("Andrew", loaded.getName());
    assertEquals(12345, loaded.getId());
  }

  @Test
  public void saveThreeTest() throws IOException {
    Collection<Student> saveds = Arrays.asList(new Student("Andrew", 1),
                                               new Student("AC", 2),
                                               new Student("Mom", 3));
    persister.save(DEFAULT_CONTEXT, serializer, saveds);

    Collection<Student> loadeds = persister.load(DEFAULT_CONTEXT, serializer);
    assertEquals(3, loadeds.size());

    Student[] loadedsArray = loadeds.toArray(new Student[0]);
    assertEquals("Andrew", loadedsArray[0].getName());
    assertEquals(1, loadedsArray[0].getId());
    assertEquals("AC", loadedsArray[1].getName());
    assertEquals(2, loadedsArray[1].getId());
    assertEquals("Mom", loadedsArray[2].getName());
    assertEquals(3, loadedsArray[2].getId());
  }

  @Test
  public void doubleContextTest() throws IOException {
    Collection<Student> savedA = Arrays.asList(new Student("Andrew", 1));
    Collection<Student> savedB = Arrays.asList(new Student("AC", 2), new Student("Mom", 3));
    persister.save(CONTEXT_A, serializer, savedA);
    persister.save(CONTEXT_B, serializer, savedB);

    Collection<Student> loadedsA = persister.load(CONTEXT_A, serializer);
    Collection<Student> loadedsB = persister.load(CONTEXT_B, serializer);
    assertEquals(1, loadedsA.size());
    assertEquals(2, loadedsB.size());

    Student[] loadedsArrayA = loadedsA.toArray(new Student[0]);
    Student[] loadedsArrayB = loadedsB.toArray(new Student[0]);
    assertEquals("Andrew", loadedsArrayA[0].getName());
    assertEquals(1, loadedsArrayA[0].getId());
    assertEquals("AC", loadedsArrayB[0].getName());
    assertEquals(2, loadedsArrayB[0].getId());
    assertEquals("Mom", loadedsArrayB[1].getName());
    assertEquals(3, loadedsArrayB[1].getId());
  }

  @Test
  public void clearTest() throws IOException {
    Student saved = new Student("Andrew", 12345);
    persister.save(DEFAULT_CONTEXT, serializer, Collections.singleton(saved));
    assertTrue(persister.exists(DEFAULT_CONTEXT));
    persister.clear(DEFAULT_CONTEXT);
    assertFalse(persister.exists(DEFAULT_CONTEXT));
  }
}
