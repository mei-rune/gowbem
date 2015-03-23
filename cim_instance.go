package gowbem

/**
 * <code>CIMElement</code> is an abstract base class that represents a CIM
 * Element as defined by the Distributed Management Task Force (<a
 * href=http://www.dmtf.org>DMTF</a>) CIM Infrastructure Specification (<a
 * href=http://www.dmtf.org/standards/published_documents/DSP0004V2.3_final.pdf
 * >DSP004</a>).
 */
type CIMElement interface {

	/**
	 * Returns a string representing the name of a CIM element instance.
	 *
	 * @return The name of this CIM element.
	 */
	GetName() string
}

/**
 * <code>CIMTypedElement</code> is an abstract class that represents a CIM
 * element that contains just the data type, but no value.
 */
type CIMTypedElement interface {
	CIMElement

	/**
	 * Returns the <code>CIMDataType</code> for this CIM Element.
	 *
	 * @return <code>CIMDataType</code> of this CIM element.
	 */
	GetType() CIMType
}

/**
 * <code>CIMValuedElement</code> is a base class used by any element that
 * contains a name, type and value.
 */
type CIMValuedElement interface {
	CIMTypedElement

	/**
	 * Returns the value for this CIM Element.
	 */
	GetValue() interface{}
}

type CIMProperty interface {
	CIMValuedElement

	/**
	 * Returns the class in which this property was defined or overridden.
	 *
	 * @return Name of class where this property was defined.
	 */
	GetOriginClass() string

	/**
	 * Convenience method for determining if this property is a Key.
	 *
	 * @return <code>true</code> if this property is a key.
	 */
	IsKey() bool

	/**
	 * Determines if this property is Propagated. When this property is part of
	 * a class, this value designates that the class origin value is the same as
	 * the class name.
	 *
	 * @return <code>true</code> if this property is propagated.
	 */
	IsPropagated() bool
}

/**
 * This class represents a CIM instance as defined by the Distributed Management
 * Task Force (<a href=http://www.dmtf.org>DMTF</a>) CIM Infrastructure
 * Specification (<a
 * href=http://www.dmtf.org/standards/published_documents/DSP0004V2.3_final.pdf
 * >DSP004</a>).
 */
type CIMInstance interface {

	/**
	 * Get the name of the class that instantiates this CIM instance.
	 *
	 * @return Name of class that instantiates this CIM instance.
	 */
	GetClassName() string

	/**
	 * Retrieve an array of the properties for this instance.
	 *
	 * @return An array of the CIM properties for this instance.
	 */
	GetProperties() []CIMProperty

	/**
	 * Get a class property by index.
	 *
	 * @param pIndex
	 *            The index of the class property to retrieve.
	 * @return The <code>CIMProperty</code> at the specified index.
	 * @throws ArrayIndexOutOfBoundsException
	 */
	GetPropertyByIndex(index int) CIMProperty

	/**
	 * Returns the specified property.
	 *
	 * @param pName
	 *            The text string for the name of the property.
	 * @return The property requested or <code>null</code> if the property does
	 *         not exist.
	 */
	GetPropertyByName(name string) CIMProperty

	/**
	 * Returns the specified <code>CIMProperty</code>.
	 *
	 * @param pName
	 *            The string name of the property to get.
	 * @param pOriginClass
	 *            (Optional) The string name of the class in which the property
	 *            was defined.
	 * @return <code>null</code> if the property does not exist, otherwise
	 *         returns the CIM property.
	 */
	GetPropertyByNameAndOrigin(name, originClass string) CIMProperty

	/**
	 * Get the number of properties defined in this <code>CIMInstance</code>.
	 *
	 * @return The number of properties defined in the <code>CIMInstance</code>.
	 */
	GetPropertyCount() int
}

type CIMKeyBinding interface {
	GetName() string
	GetType() CIMType
	GetValue() interface{}
}

type CIMKeyBindings interface {
	Len() int
	Get(idx int) CIMKeyBinding
}

type CIMInstanceName interface {
	GetClassName() string
	GetKeyBindings() CIMKeyBindings
	String() string
}

/**
 * This class represents the CIM Object Path as defined by the Distributed
 * Management Task Force (<a href=http://www.dmtf.org>DMTF</a>) CIM
 * Infrastructure Specification (<a href=
 * "http://dmtf.org/sites/default/files/standards/documents/DSP0004_2.7.0.pdf"
 * >DSP004</a>). In order to uniquely identify a given object, a CIM object path
 * includes the host, namespace, object name and keys (if the object is an
 * instance).<br>
 * <br>
 * For example, the object path:<br>
 * <br>
 * <code>
 * http://myserver/root/cimv2:My_ComputerSystem.Name=mycomputer,
 * CreationClassName=My_ComputerSystem
 * </code><br>
 * <br>
 * has two parts:<br>
 * <br>
 * <ul type="disc"> <li>Namespace Path</li>
 * <code>http://myserver/root/cimv2</code><br>
 * JSR48 defines the namespace path to include the scheme, host, port (optional)
 * and namespace<br>
 * The example specifies the <code>"root/cimv2"</code> namespace on the host
 * <code>myserver</code>.</li> <li>Model Path</li>
 * <code>My_ComputerSystem.Name=mycomputer,CreationClassName=My_ComputerSystem
 * 		</code><br>
 * DSP0004 defines the model path for a class or qualifier type as the name of
 * the class/qualifier type<br>
 * DSP0004 defines the model path for an instance as the class
 * name.(key=value),*<br>
 * The example specifies an instance for the class
 * <code>My_ComputerSystem</code> which is uniquely identified by two key
 * properties and values: <ul type="disc"> <li><code>Name=mycomputer</code></li>
 * <li>
 * <code>CreationClassName=My_ComputerSystem</code></li> </ul> </ul>
 */
type CIMObjectPath interface {

	/**
	 * Gets the host.
	 */
	GetHost() string

	/**
	 * Gets a key property by name.
	 *
	 * @param pName
	 *            The name of the key property to retrieve.
	 * @return The <code>CIMProperty</code> with the given name, or
	 *         <code>null</code> if it is not found.
	 */
	GetKey(name string) CIMValuedElement

	/**
	 * Gets all key properties.
	 *
	 * @return The container of key properties.
	 */
	GetKeys() map[string]CIMValuedElement

	/**
	 * Gets the namespace.
	 *
	 * @return The name of the namespace.
	 */
	GetNamespace() string

	/**
	 * Gets the object name. Depending on the type of CIM element referenced,
	 * this may be either a class name or a qualifier type name.
	 *
	 * @return The name of this CIM element.
	 */
	GetObjectName() string

	/**
	 * Gets the the port on the host to which the connection was established.
	 *
	 * @return The port on the host.
	 */
	GetPort() string

	/**
	 * Get the connection scheme.
	 *
	 * @return The connection scheme (e.g. http, https,...)
	 */
	GetScheme() string
}

type CIMInstanceWithName interface {
	GetName() CIMInstanceName
	GetInstance() CIMInstance
}
